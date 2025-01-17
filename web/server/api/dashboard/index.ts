export default defineEventHandler(async (event) => {
	// make sure the user is logged in
	// This will throw a 401 error if the request doesn't come from a valid user session
	const { user } = await requireUserSession(event);

	// TODO: Fetch some stats based on the user

	return {};
});
