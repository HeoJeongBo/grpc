import { TransportProvider } from "@connectrpc/connect-query";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { createRouter, RouterProvider } from "@tanstack/react-router";
import React from "react";
import ReactDOM from "react-dom/client";
import { routeTree } from "./routeTree.gen";
import "./index.css";
import { AuthProvider, useAuth } from "./context/auth-context";
import { transport } from "./lib/connect-client";

// biome-ignore lint/style/noNonNullAssertion: initialize auth with undefined
const router = createRouter({ routeTree, context: { auth: undefined! } });

const queryClient = new QueryClient();
declare module "@tanstack/react-router" {
	interface Register {
		router: typeof router;
	}
}

function App() {
	const auth = useAuth();

	return <RouterProvider router={router} context={{ auth }} />;
}

// biome-ignore lint/style/noNonNullAssertion: root is always exist
ReactDOM.createRoot(document.getElementById("root")!).render(
	<React.StrictMode>
		<AuthProvider>
			<TransportProvider transport={transport}>
				<QueryClientProvider client={queryClient}>
					<App />
				</QueryClientProvider>
			</TransportProvider>
		</AuthProvider>
	</React.StrictMode>,
);
