import { ButtonIcon, ButtonRoot, ButtonText } from "./button";

export const Button = Object.assign(ButtonRoot, {
	Icon: ButtonIcon,
	Text: ButtonText,
});

export type { ButtonProps } from "./button";
export { buttonVariants } from "./button";
