
export default {
  async fetch(request) {
    const target = "https://my.koom.pp.ua";

    const upgrade = request.headers.get("Upgrade") || "";
    if (upgrade.toLowerCase() === "websocket") {
      const backend = await fetch(target, {
        method: request.method,
        headers: request.headers,
      });
      return backend;
    }

    return fetch(target + new URL(request.url).pathname, {
      method: request.method,
      headers: request.headers,
      body: request.body,
    });
  }
}
