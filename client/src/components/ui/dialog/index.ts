import {
	DialogClose,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogOverlay,
	DialogPortal,
	DialogRoot,
	DialogTitle,
	DialogTrigger,
} from "./dialog";

export const Dialog = Object.assign(DialogRoot, {
	Trigger: DialogTrigger,
	Portal: DialogPortal,
	Overlay: DialogOverlay,
	Content: DialogContent,
	Header: DialogHeader,
	Footer: DialogFooter,
	Title: DialogTitle,
	Description: DialogDescription,
	Close: DialogClose,
});
