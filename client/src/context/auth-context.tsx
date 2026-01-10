import type React from "react";
import { createContext, useCallback, useContext, useMemo } from "react";
import { useAuthStore } from "../store/auth-store";

export interface AuthContext {
	accessToken: string | null;
	refreshToken: string | null;
	isAuthenticated: boolean;
	login: (accessToken: string, refreshToken: string) => void;
	logout: () => void;
	updateAccessToken: (accessToken: string) => void;
}

const AuthContext = createContext<AuthContext | null>(null);

export interface AuthProviderProps {
	children: React.ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
	const {
		accessToken,
		refreshToken,
		isAuthenticated,
		setTokens,
		clearTokens,
		updateAccessToken: storeUpdateAccessToken,
	} = useAuthStore();

	const login = useCallback(
		(accessToken: string, refreshToken: string) => {
			setTokens(accessToken, refreshToken);
		},
		[setTokens],
	);

	const logout = useCallback(() => {
		clearTokens();
	}, [clearTokens]);

	const updateAccessToken = useCallback(
		(accessToken: string) => {
			storeUpdateAccessToken(accessToken);
		},
		[storeUpdateAccessToken],
	);

	const value = useMemo(
		() => ({
			accessToken,
			refreshToken,
			isAuthenticated,
			login,
			logout,
			updateAccessToken,
		}),
		[
			accessToken,
			refreshToken,
			isAuthenticated,
			login,
			logout,
			updateAccessToken,
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
