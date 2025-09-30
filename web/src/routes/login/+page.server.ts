import { type User } from '$lib/user';
import type { Actions } from './$types';

export const actions: Actions = {
    login: async ({ request, fetch }) => {
        const data = await request.formData();
        const params: User = {
            name: data.get("name") as string,
            password: data.get("password") as string
        };

        // Try sign user in
        const form = new FormData();
        form.append('name', params.name);
        form.append('password', params.password);

        // Call your POST endpoint
        try {
            const response = await fetch('/api/login', {
                method: 'POST',
                body: form
            });
            if (!response.ok) {
                const err : {errMsg: string} = await response.json()
                throw new Error(err.errMsg)
            }
        } catch (err) {
            if (err instanceof Error) {
            return {
                errMsg: err.message,
                user: {
                    name: params.name,
                    password: params.password
                }
            }
            }
        }
        
    }
};