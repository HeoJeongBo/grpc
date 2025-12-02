import { createPromiseClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { ItemService } from "@/proto-generated/item/v1/item_connect";

const transport = createConnectTransport({
  baseUrl: "http://localhost:8080",
});

export const itemClient = createPromiseClient(ItemService, transport);
