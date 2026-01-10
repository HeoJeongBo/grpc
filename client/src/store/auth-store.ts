import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

interface AuthState {
	accessToken: string | null;
	refreshToken: string | null;
	isAuthenticated: boolean;
}

interface AuthActions {
	setTokens: (accessToken: string, refreshToken: string) => void;
	clearTokens: () => void;
	updateAccessToken: (accessToken: string) => void;
}

type AuthStore = AuthState & AuthActions;

const initialState: AuthState = {
	accessToken: null,
	refreshToken: null,
	isAuthenticated: false,
};

export const useAuthStore = create<AuthStore>()(
	persist(
		(set) => ({
			...initialState,
			setTokens: (accessToken, refreshToken) => {
				set({
					accessToken,
					refreshToken,
					isAuthenticated: true,
				});
			},
			clearTokens: () => set(initialState),
			updateAccessToken: (accessToken) =>
				set((state) => ({
					...state,
					accessToken,
				})),
		}),
		{
			name: "auth-storage",
			storage: createJSONStorage(() => localStorage),
		},
	),
);
