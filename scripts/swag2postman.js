const fs = require('fs');
const path = require('path');

const collectionPath = path.join(__dirname, '../.postman/Dietician Backend API.postman_collection.json');
let collection;
try {
    collection = JSON.parse(fs.readFileSync(collectionPath, 'utf-8'));
} catch (e) {
    console.error("Could not read existing Postman collection:", e);
    process.exit(1);
}

// Build a map of existing requests to preserve bodies and scripts
const existingRequests = new Map();
function mapExistingItems(itemArray) {
    for (const item of itemArray) {
        if (item.item) {
            mapExistingItems(item.item);
        } else if (item.request) {
            // Key based on method + url path
            const method = (item.request.method || 'GET').toUpperCase();
            let urlPath = '';
            if (item.request.url && item.request.url.path) {
                urlPath = '/' + item.request.url.path.join('/');
            } else if (typeof item.request.url === 'string') {
                const parts = item.request.url.split('}}');
                urlPath = parts.length > 1 ? parts[1] : parts[0];
            }
            const key = `${method} ${urlPath}`;
            existingRequests.set(key, item);
        }
    }
}
mapExistingItems(collection.item);

const servicesDir = path.join(__dirname, '../services');
const services = fs.readdirSync(servicesDir).filter(f => fs.statSync(path.join(servicesDir, f)).isDirectory());

const newItems = [];

for (const service of services) {
    const swaggerPath = path.join(servicesDir, service, 'docs/swagger.json');
    if (!fs.existsSync(swaggerPath)) continue;

    const swagger = JSON.parse(fs.readFileSync(swaggerPath, 'utf-8'));
    
    // Determine service name and variable
    // e.g. "account-service" -> "Account Service", variable "{{account_service}}"
    const serviceName = service.split('-').map(w => w.charAt(0).toUpperCase() + w.slice(1)).join(' ');
    const serviceVar = `{{${service.replace(/-/g, '_')}}}`;
    
    const serviceFolder = {
        name: serviceName,
        item: []
    };
    
    const tagsMap = new Map();

    for (const [apiPath, methods] of Object.entries(swagger.paths || {})) {
        for (const [method, details] of Object.entries(methods)) {
            const tag = (details.tags && details.tags.length > 0) ? details.tags[0] : 'Default';
            
            if (!tagsMap.has(tag)) {
                tagsMap.set(tag, {
                    name: tag,
                    item: []
                });
            }
            
            // Build key to look up existing
            // apiPath usually is like /api/v1/auth/register
            const key = `${method.toUpperCase()} ${apiPath}`;
            const existing = existingRequests.get(key);

            // Construct URL parts
            const urlParts = apiPath.split('/').filter(p => p.length > 0);
            
            // Clean up path parameters e.g., {id} -> :id (postman format)
            const postmanUrlParts = urlParts.map(p => p.startsWith('{') && p.endsWith('}') ? ':' + p.slice(1, -1) : p);
            const pathVariables = [];
            urlParts.forEach(p => {
                if (p.startsWith('{') && p.endsWith('}')) {
                    pathVariables.push({
                        key: p.slice(1, -1),
                        value: ""
                    });
                }
            });

            // Reconstruct Postman item
            const postmanItem = {
                name: details.summary || apiPath,
                request: {
                    method: method.toUpperCase(),
                    header: [],
                    url: {
                        raw: `{{base_url}}${serviceVar}${apiPath.replace(/{([^}]+)}/g, ':$1')}`,
                        host: [
                            `{{base_url}}${serviceVar}`
                        ],
                        path: postmanUrlParts
                    }
                },
                response: []
            };

            if (pathVariables.length > 0) {
                postmanItem.request.url.variable = pathVariables;
            }
            
            // Query parameters
            const queryParams = (details.parameters || []).filter(p => p.in === 'query');
            if (queryParams.length > 0) {
                postmanItem.request.url.query = queryParams.map(q => ({
                    key: q.name,
                    value: "",
                    description: q.description || ""
                }));
            }

            // Headers
            postmanItem.request.header = [
                {
                    key: "Accept",
                    value: "application/json"
                }
            ];

            // If requires auth
            if (details.security && details.security.length > 0) {
                postmanItem.request.auth = {
                    type: "bearer",
                    bearer: [
                        {
                            key: "token",
                            value: "{{access_token}}",
                            type: "string"
                        }
                    ]
                };
            }

            // Merge with existing
            if (existing) {
                if (existing.event) {
                    postmanItem.event = existing.event;
                }
                if (existing.request.body) {
                    postmanItem.request.body = existing.request.body;
                    if (!postmanItem.request.header.some(h => h.key.toLowerCase() === 'content-type')) {
                        postmanItem.request.header.push({
                            key: "Content-Type",
                            value: "application/json"
                        });
                    }
                }
            } else if (method.toLowerCase() === 'post' || method.toLowerCase() === 'put' || method.toLowerCase() === 'patch') {
                postmanItem.request.body = {
                    mode: "raw",
                    raw: "{}"
                };
                postmanItem.request.header.push({
                    key: "Content-Type",
                    value: "application/json"
                });
            }

            tagsMap.get(tag).item.push(postmanItem);
        }
    }
    
    // Add tags to service folder
    for (const tagFolder of tagsMap.values()) {
        serviceFolder.item.push(tagFolder);
    }
    
    if (serviceFolder.item.length > 0) {
        newItems.push(serviceFolder);
    }
}

collection.item = newItems;

fs.writeFileSync(collectionPath, JSON.stringify(collection, null, 2));
console.log('Successfully updated Postman collection from swagger files.');
