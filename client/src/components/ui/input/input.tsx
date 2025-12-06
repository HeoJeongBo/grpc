import * as LabelPrimitive from "@radix-ui/react-label";
import { cva, type VariantProps } from "class-variance-authority";
import {
	type ComponentProps,
	type ComponentPropsWithoutRef,
	type HTMLAttributes,
	useId,
} from "react";

import { cn } from "@/lib/utils";

const inputVariants = cva(
	"flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-base shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 md:text-sm",
);

const labelVariants = cva(
	"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
);

type InputFieldProps = {
	label?: string;
	error?: string;
	helperText?: string;
} & ComponentProps<"input">;

const InputRoot = ({ className, ...props }: HTMLAttributes<HTMLDivElement>) => {
	return <div className={cn("space-y-2", className)} {...props} />;
};

type InputLabelProps = ComponentPropsWithoutRef<typeof LabelPrimitive.Root> &
	VariantProps<typeof labelVariants>;

const InputLabel = ({ className, ...props }: InputLabelProps) => (
	<LabelPrimitive.Root className={cn(labelVariants(), className)} {...props} />
);

const InputControl = ({
	className,
	type,
	...props
}: ComponentProps<"input">) => {
	return (
		<input type={type} className={cn(inputVariants(), className)} {...props} />
	);
};

const InputError = ({
	className,
	children,
	...props
}: HTMLAttributes<HTMLParagraphElement>) => {
	if (!children) return null;

	return (
		<p className={cn("text-sm text-destructive", className)} {...props}>
			{children}
		</p>
	);
};

const InputHelper = ({
	className,
	children,
	...props
}: HTMLAttributes<HTMLParagraphElement>) => {
	if (!children) return null;

	return (
		<p className={cn("text-sm text-muted-foreground", className)} {...props}>
			{children}
		</p>
	);
};

// Convenience component for simple inputs
const InputField = ({
	label,
	error,
	helperText,
	className,
	id,
	...props
}: InputFieldProps) => {
	const generatedId = useId();
	const inputId = id || generatedId;

	return (
		<InputRoot className={className}>
			{label && <InputLabel htmlFor={inputId}>{label}</InputLabel>}
			<InputControl id={inputId} {...props} />
			<InputError>{error}</InputError>
			<InputHelper>{helperText}</InputHelper>
		</InputRoot>
	);
};

export {
	InputRoot,
	InputLabel,
	InputControl,
	InputError,
	InputHelper,
	InputField,
};
