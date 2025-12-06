import {
	InputControl,
	InputError,
	InputField,
	InputHelper,
	InputLabel,
	InputRoot,
} from "./input";

export const Input = Object.assign(InputRoot, {
	Label: InputLabel,
	Control: InputControl,
	Error: InputError,
	Helper: InputHelper,
	Field: InputField,
});

export { InputField };
