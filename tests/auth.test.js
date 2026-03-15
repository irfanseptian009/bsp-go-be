const { createApiClient } = require('./helpers/apiClient.cjs');

describe('Auth API', () => {
  const api = createApiClient();

  test('login admin berhasil', async () => {
    const response = await api.post('/api/auth/login', {
      email: 'admin@fims.com',
      password: 'admin123',
    });

    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('accessToken');
    expect(response.data).toHaveProperty('user');
    expect(response.data.user.role).toBe('ADMIN');
  });

  test('login gagal untuk password salah', async () => {
    await expect(
      api.post('/api/auth/login', {
        email: 'admin@fims.com',
        password: 'salah-total',
      }),
    ).rejects.toMatchObject({
      response: {
        status: 401,
      },
    });
  });
});
