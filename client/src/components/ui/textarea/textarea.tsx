import * as LabelPrimitive from "@radix-ui/react-label";
import { cva, type VariantProps } from "class-variance-authority";
import {
	type ComponentProps,
	type ComponentPropsWithoutRef,
	type HTMLAttributes,
	useId,
} from "react";

import { cn } from "@/lib/utils";

const textareaVariants = cva(
	"flex min-h-[60px] w-full rounded-md border border-input bg-transparent px-3 py-2 text-base shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 md:text-sm",
);

const labelVariants = cva(
	"text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
);

type TextareaFieldProps = {
	label?: string;
	error?: string;
	helperText?: string;
} & ComponentProps<"textarea">;

const TextareaRoot = ({
	className,
	...props
}: HTMLAttributes<HTMLDivElement>) => {
	return <div className={cn("space-y-2", className)} {...props} />;
};

type TextareaLabelProps = ComponentPropsWithoutRef<typeof LabelPrimitive.Root> &
	VariantProps<typeof labelVariants>;

const TextareaLabel = ({ className, ...props }: TextareaLabelProps) => (
	<LabelPrimitive.Root className={cn(labelVariants(), className)} {...props} />
);

const TextareaControl = ({
	className,
	...props
}: ComponentProps<"textarea">) => {
	return <textarea className={cn(textareaVariants(), className)} {...props} />;
};

const TextareaError = ({
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

const TextareaHelper = ({
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

// Convenience component for simple textareas
const TextareaField = ({
	label,
	error,
	helperText,
	className,
	id,
	...props
}: TextareaFieldProps) => {
	const generatedId = useId();
	const textareaId = id || generatedId;

	return (
		<TextareaRoot className={className}>
			{label && <TextareaLabel htmlFor={textareaId}>{label}</TextareaLabel>}
			<TextareaControl id={textareaId} {...props} />
			<TextareaError>{error}</TextareaError>
			<TextareaHelper>{helperText}</TextareaHelper>
		</TextareaRoot>
	);
};

export {
	TextareaRoot,
	TextareaLabel,
	TextareaControl,
	TextareaError,
	TextareaHelper,
	TextareaField,
};
