import { Container, getContainer, getRandom } from "@cloudflare/containers";
import { Hono } from "hono";

export class MyContainer extends Container<Env> {
	defaultPort = 8080;
	sleepAfter = "2m";

	envVars = {
		MESSAGE: "Proxy iniciado dentro do container!"
	};

	override onStart() {
		console.log("Container iniciado");
	}

	override onStop() {
		console.log("Container finalizado");
	}

	override onError(error: unknown) {
		console.log("Erro no container:", error);
	}
}

const app = new Hono<{ Bindings: Env }>();

app.get("/", (c) =>
	c.text(
		"Endpoints disponíveis:\n" +
			"GET /container/<ID>\n" +
			"GET /lb\n" +
			"GET /singleton\n"
	),
);

// 1) Criar containers por ID
app.get("/container/:id", async (c) => {
	const id = c.req.param("id");

	const containerId = c.env.MY_CONTAINER.idFromName(`/container/${id}`);
	const container = c.env.MY_CONTAINER.get(containerId);

	return await container.fetch(c.req.raw);
});

// 2) Load balance entre containers
app.get("/lb", async (c) => {
	const container = await getRandom(c.env.MY_CONTAINER, 3);
	return await container.fetch(c.req.raw);
});

// 3) Instância única
app.get("/singleton", async (c) => {
	const container = getContainer(c.env.MY_CONTAINER);
	return await container.fetch(c.req.raw);
});

export default app;
