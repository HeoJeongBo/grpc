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
import { register } from "@/proto-generated/auth/auth_service-AuthService_connectquery";

const signUpSchema = z
	.object({
		email: z
			.string()
			.min(1, "이메일을 입력하세요")
			.email("올바른 이메일 형식이 아닙니다"),
		password: z.string().min(6, "비밀번호는 최소 6자 이상이어야 합니다"),
		confirmPassword: z.string().min(6, "비밀번호를 확인해주세요"),
		name: z.string().min(1, "이름을 입력하세요"),
	})
	.refine((data) => data.password === data.confirmPassword, {
		message: "비밀번호가 일치하지 않습니다",
		path: ["confirmPassword"],
	});

type SignUpData = z.infer<typeof signUpSchema>;

export const Route = createFileRoute("/(non-auth)/sign-up")({
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

	const registerMutation = useMutation(register, {
		onSuccess: (data) => {
			if (data.user && data.tokens) {
				login(data.user, data.tokens.accessToken, data.tokens.refreshToken);
				toast.success("회원가입 성공!");
				navigate({ to: "/" });
			}
		},
		onError: (error) => {
			toast.error(`회원가입 실패: ${error.message}`);
		},
	});

	const form = useForm({
		defaultValues: {
			email: "",
			password: "",
			confirmPassword: "",
			name: "",
		} as SignUpData,
		onSubmit: async ({ value }) => {
			const result = signUpSchema.safeParse(value);
			if (!result.success) {
				toast.error("입력값을 확인해주세요");
				return;
			}
			registerMutation.mutate({
				email: result.data.email,
				password: result.data.password,
				name: result.data.name,
			});
		},
	});

	return (
		<div className="flex min-h-screen items-center justify-center">
			<div className="w-full max-w-md space-y-8 rounded-lg border p-8">
				<div className="text-center">
					<h2 className="text-3xl font-bold">회원가입</h2>
					<p className="mt-2 text-sm text-muted-foreground">
						새 계정을 만드세요
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
						name="name"
						validators={{
							onChange: ({ value }) => {
								const result = signUpSchema.shape.name.safeParse(value);
								if (!result.success) {
									return result.error.issues[0]?.message;
								}
								return undefined;
							},
						}}
					>
						{(field) => (
							<InputRoot>
								<InputLabel htmlFor="name">이름</InputLabel>
								<InputControl
									id="name"
									type="text"
									value={field.state.value}
									onBlur={field.handleBlur}
									onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
										field.handleChange(e.target.value)
									}
									placeholder="홍길동"
								/>
								<InputError>{field.state.meta.errors[0]}</InputError>
							</InputRoot>
						)}
					</form.Field>

					<form.Field
						name="email"
						validators={{
							onChange: ({ value }) => {
								const result = signUpSchema.shape.email.safeParse(value);
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
								const result = signUpSchema.shape.password.safeParse(value);
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

					<form.Field
						name="confirmPassword"
						validators={{
							onChangeListenTo: ["password"],
							onChange: ({ value, fieldApi }) => {
								const password = fieldApi.form.getFieldValue("password");
								if (value !== password) {
									return "비밀번호가 일치하지 않습니다";
								}
								return undefined;
							},
						}}
					>
						{(field) => (
							<InputRoot>
								<InputLabel htmlFor="confirmPassword">비밀번호 확인</InputLabel>
								<InputControl
									id="confirmPassword"
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
						disabled={registerMutation.isPending}
					>
						{registerMutation.isPending ? "가입 중..." : "회원가입"}
					</ButtonRoot>

					<div className="text-center text-sm">
						<span className="text-muted-foreground">
							이미 계정이 있으신가요?{" "}
						</span>
						<button
							type="button"
							onClick={() => navigate({ to: "/sign-in" })}
							className="text-primary hover:underline"
						>
							로그인
						</button>
					</div>
				</form>
			</div>
		</div>
	);
}
