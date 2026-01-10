import { useMutation } from "@connectrpc/connect-query";
import { useForm } from "@tanstack/react-form";
import { createFileRoute, redirect, useNavigate } from "@tanstack/react-router";
import { toast } from "sonner";
import { z } from "zod";
import { ButtonRoot } from "@/components/ui/button/button";
import {
	InputControl,
	InputError,
	InputLabel,
	InputRoot,
} from "@/components/ui/input/input";
import { useAuth } from "@/context/auth-context";
import { login as loginMethod } from "@/proto-generated/auth/auth_service-AuthService_connectquery";

const loginSchema = z.object({
	email: z
		.string()
		.min(1, "이메일을 입력하세요")
		.email("올바른 이메일 형식이 아닙니다"),
	password: z.string().min(6, "비밀번호는 최소 6자 이상이어야 합니다"),
});

type LoginData = z.infer<typeof loginSchema>;

export const Route = createFileRoute("/(non-auth)/sign-in")({
	component: RouteComponent,
	beforeLoad: ({ context }) => {
		if (context.auth.isAuthenticated) {
			throw redirect({
				to: "/",
			});
		}
	},
});

function RouteComponent() {
	const navigate = useNavigate();
	const { login } = useAuth();

	const loginMutation = useMutation(loginMethod, {
		onSuccess: (data) => {
			if (data.user && data.tokens) {
				login(
					{
						email: data.user.email,
						id: data.user.id,
						name: data.user.name,
					},
					data.tokens.accessToken,
					data.tokens.refreshToken,
				);
				toast.success("로그인 성공!");
				navigate({ to: "/" });
			}
		},
		onError: (error) => {
			toast.error(`로그인 실패: ${error.message}`);
		},
	});

	const form = useForm({
		defaultValues: {
			email: "",
			password: "",
		} as LoginData,
		onSubmit: async ({ value }) => {
			const result = loginSchema.safeParse(value);
			if (!result.success) {
				toast.error("입력값을 확인해주세요");
				return;
			}
			loginMutation.mutate(result.data);
		},
	});

	return (
		<div className="flex min-h-screen items-center justify-center">
			<div className="w-full max-w-md space-y-8 rounded-lg border p-8">
				<div className="text-center">
					<h2 className="text-3xl font-bold">로그인</h2>
					<p className="mt-2 text-sm text-muted-foreground">
						계정에 로그인하세요
					</p>
				</div>

				<form
					onSubmit={(e) => {
						e.preventDefault();
						e.stopPropagation();
						form.handleSubmit();
					}}
					className="space-y-6"
				>
					<form.Field
						name="email"
						validators={{
							onChange: ({ value }) => {
								const result = loginSchema.shape.email.safeParse(value);
								if (!result.success) {
									return result.error.issues[0]?.message;
								}
								return undefined;
							},
						}}
					>
						{(field) => (
							<InputRoot>
								<InputLabel htmlFor="email">이메일</InputLabel>
								<InputControl
									id="email"
									type="email"
									value={field.state.value}
									onBlur={field.handleBlur}
									onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
										field.handleChange(e.target.value)
									}
									placeholder="email@example.com"
								/>
								<InputError>{field.state.meta.errors[0]}</InputError>
							</InputRoot>
						)}
					</form.Field>

					<form.Field
						name="password"
						validators={{
							onChange: ({ value }) => {
								const result = loginSchema.shape.password.safeParse(value);
								if (!result.success) {
									return result.error.issues[0]?.message;
								}
								return undefined;
							},
						}}
					>
						{(field) => (
							<InputRoot>
								<InputLabel htmlFor="password">비밀번호</InputLabel>
								<InputControl
									id="password"
									type="password"
									value={field.state.value}
									onBlur={field.handleBlur}
									onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
										field.handleChange(e.target.value)
									}
								/>
								<InputError>{field.state.meta.errors[0]}</InputError>
							</InputRoot>
						)}
					</form.Field>

					<ButtonRoot
						type="submit"
						className="w-full"
						disabled={loginMutation.isPending}
					>
						{loginMutation.isPending ? "로그인 중..." : "로그인"}
					</ButtonRoot>

					<div className="text-center text-sm">
						<span className="text-muted-foreground">계정이 없으신가요? </span>
						<button
							type="button"
							onClick={() => navigate({ to: "/sign-up" })}
							className="text-primary hover:underline"
						>
							회원가입
						</button>
					</div>
				</form>
			</div>
		</div>
	);
}
