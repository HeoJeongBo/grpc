import type React from "react";
import { createContext, useCallback, useContext, useMemo } from "react";
import type { User } from "@/proto-generated/user/user_pb";
import { useAuthStore } from "../store/auth-store";

export interface AuthContext {
	user: User | null;
	accessToken: string | null;
	refreshToken: string | null;
	isAuthenticated: boolean;
	login: (user: User, accessToken: string, refreshToken: string) => void;
	logout: () => void;
	updateAccessToken: (accessToken: string) => void;
	updateUser: (user: User) => void;
}

const AuthContext = createContext<AuthContext | null>(null);

export interface AuthProviderProps {
	children: React.ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
	const {
		user,
		accessToken,
		refreshToken,
		isAuthenticated,
		setAuth,
		clearAuth,
		updateAccessToken: storeUpdateAccessToken,
		setUser,
	} = useAuthStore();

	const login = useCallback(
		(user: User, accessToken: string, refreshToken: string) => {
			setAuth(user, accessToken, refreshToken);
		},
		[setAuth],
	);

	const logout = useCallback(() => {
		clearAuth();
	}, [clearAuth]);

	const updateAccessToken = useCallback(
		(accessToken: string) => {
			storeUpdateAccessToken(accessToken);
		},
		[storeUpdateAccessToken],
	);

	const updateUser = useCallback(
		(user: User) => {
			setUser(user);
		},
		[setUser],
	);

	const value = useMemo(
		() => ({
			user,
			accessToken,
			refreshToken,
			isAuthenticated,
			login,
			logout,
			updateAccessToken,
			updateUser,
		}),
		[
			user,
			accessToken,
			refreshToken,
			isAuthenticated,
			login,
			logout,
			updateAccessToken,
			updateUser,
		],
	);

	return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
	const context = useContext(AuthContext);
	if (!context) {
		throw new Error("useAuth must be used within an AuthProvider");
	}
	return context;
}
