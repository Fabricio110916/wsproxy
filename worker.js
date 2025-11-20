export default {
  async fetch(request) {
    const url = new URL(request.url);

    // backend base
    const target = "https://my.koom.pp.ua";

    // destino completo preservando path + query
    const dest = target + url.pathname + url.search;

    // recria a request original apontando para o backend real
    const newReq = new Request(dest, {
      method: request.method,
      headers: request.headers,
      body: request.body,
      duplex: "half"   // necessário para WebSocket + body streaming
    });

    // encaminha tudo (HTTP normal + WebSocket)
    return fetch(newReq);
  }
}
