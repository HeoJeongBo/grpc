import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/(non-auth)/sign-in")({
	component: RouteComponent,
	beforeLoad: ({ context }) => {},
});

function RouteComponent() {
	return <div>Hello "/(non-auth)/sign-in"!</div>;
}
