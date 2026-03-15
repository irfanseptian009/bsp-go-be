const { createApiClient } = require('./helpers/apiClient.cjs');

async function getAdminToken() {
  const api = createApiClient();
  const response = await api.post('/api/auth/login', {
    email: 'admin@fims.com',
    password: 'admin123',
  });

  return response.data.accessToken;
}

describe('Occupation Types API', () => {
  let token;
  let createdId;

  beforeAll(async () => {
    token = await getAdminToken();
  });

  afterAll(async () => {
    if (!createdId) return;

    const api = createApiClient(token);
    try {
      await api.delete(`/api/occupation-types/${createdId}`);
    } catch (error) {
      // ignore cleanup errors
    }
  });

  test('admin bisa create occupation type', async () => {
    const api = createApiClient(token);
    const uniqueCode = `JEST.${Date.now()}`;

    const response = await api.post('/api/occupation-types', {
      code: uniqueCode,
      name: 'Occupation Jest Test',
      premiumRate: 1.1111,
    });

    expect(response.status).toBe(201);
    expect(response.data.id).toBeTruthy();
    expect(response.data.code).toBe(uniqueCode);

    createdId = response.data.id;
  });

  test('admin tidak bisa create dengan code duplikat', async () => {
    const api = createApiClient(token);

    const payload = {
      code: '2976.01',
      name: 'Duplicate Code',
      premiumRate: 0.9,
    };

    await expect(api.post('/api/occupation-types', payload)).rejects.toMatchObject({
      response: {
        status: 409,
      },
    });
  });

  test('customer tidak boleh create occupation type', async () => {
    const anonymousApi = createApiClient();
    const customerLogin = await anonymousApi.post('/api/auth/login', {
      email: 'customer@fims.com',
      password: 'customer123',
    });

    const customerApi = createApiClient(customerLogin.data.accessToken);

    await expect(
      customerApi.post('/api/occupation-types', {
        code: `CUS.${Date.now()}`,
        name: 'Should Fail',
        premiumRate: 1.2,
      }),
    ).rejects.toMatchObject({
      response: {
        status: 403,
      },
    });
  });
});
