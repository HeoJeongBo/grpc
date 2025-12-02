import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { ItemService } from "@/proto-generated/item/v1/item_pb";

const transport = createConnectTransport({
	baseUrl: "http://localhost:8080",
});

export const itemClient = createClient(ItemService, transport);
