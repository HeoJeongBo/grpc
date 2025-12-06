import {
	TextareaControl,
	TextareaError,
	TextareaField,
	TextareaHelper,
	TextareaLabel,
	TextareaRoot,
} from "./textarea";

export const Textarea = Object.assign(TextareaRoot, {
	Label: TextareaLabel,
	Control: TextareaControl,
	Error: TextareaError,
	Helper: TextareaHelper,
	Field: TextareaField,
});

export { TextareaField };
