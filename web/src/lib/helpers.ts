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

interface Quote {
	fact: string;
	length: number;
}

export async function fetchQuote(): Promise<string> {
	try {
		const response = await fetch('https://catfact.ninja/fact');
		if (!response.ok) {
			return "Services don't always work as expected NOT OK";
		}
		const parsedResponse: Quote = await response.json();
		return parsedResponse.fact;
	} catch (err) {
		console.error(err);
		return "Services don't always work as expected ERROR";
	}
}
