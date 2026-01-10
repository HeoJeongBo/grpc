import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";
import type { User } from "@/proto-generated/user/user_pb";

interface AuthState {
	user: User | null;
	accessToken: string | null;
	refreshToken: string | null;
	isAuthenticated: boolean;
}

interface AuthActions {
	setAuth: (user: User, accessToken: string, refreshToken: string) => void;
	clearAuth: () => void;
	updateAccessToken: (accessToken: string) => void;
	setUser: (user: User) => void;
}

type AuthStore = AuthState & AuthActions;

const initialState: AuthState = {
	user: null,
	accessToken: null,
	refreshToken: null,
	isAuthenticated: false,
};

export const useAuthStore = create<AuthStore>()(
	persist(
		(set) => ({
			...initialState,
			setAuth: (user, accessToken, refreshToken) =>
				set({
					user,
					accessToken,
					refreshToken,
					isAuthenticated: true,
				}),
			clearAuth: () => set(initialState),
			updateAccessToken: (accessToken) =>
				set((state) => ({
					...state,
					accessToken,
				})),
			setUser: (user) =>
				set((state) => ({
					...state,
					user,
				})),
		}),
		{
			name: "auth-storage",
			storage: createJSONStorage(() => localStorage),
		},
	),
);
