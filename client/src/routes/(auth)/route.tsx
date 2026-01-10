import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/(auth)")({
	component: RouteComponent,
	beforeLoad: ({ context }) => {
		if (!context.auth.isAuthenticated) {
			throw redirect({
				to: "/sign-in",
			});
		}
	},
});

function RouteComponent() {
	return <div>Hello "/(auth)"!</div>;
}
