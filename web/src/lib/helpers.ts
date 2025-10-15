import { toast } from 'svelte-sonner';

export function logError(msgToUser: string, err: Error) {
	toast.error(msgToUser);

	console.error(JSON.stringify(err));
}

export function logSuccess(msg: string) {
	toast.success(msg);
}

export function truncateString(text: string, maxLength: number): string {
	if (text.length <= maxLength) {
		return text;
	}
	return text.substring(0, maxLength) + '...';
}
