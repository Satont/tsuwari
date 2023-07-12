import { useQuery, useMutation } from '@tanstack/vue-query';
import { Profile } from '@twir/grpc/generated/api/api/auth';
import { ref } from 'vue';

import { protectedApiClient } from './twirp.js';

const profile = ref<Profile | null>(null);

export const useProfile = () =>
	useQuery<Profile | null>({
		queryKey: [`authProfile`],
		queryFn: async () => {
			const user = await getProfile();
			profile.value = user;
			return user;
		},
		retry: false,
		initialData: profile.value,
	});

export const useLogout = () => useMutation({
	mutationKey: ['authLogout'],
	mutationFn: async () => {
		await protectedApiClient.authLogout({});
		profile.value = null;
	},
	onSuccess: () => {
		profile.value = null;
		window.location.replace('/');
	},
});

export const getProfile = async () => {
	if (profile.value) return profile.value;

	const call = await protectedApiClient.authUserProfile({});
	profile.value = call.response;
	return call.response;
};
