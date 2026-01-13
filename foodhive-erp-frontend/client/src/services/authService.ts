import api from '@/lib/api';

export const authService = {
  login: async (credentials: { email: string; password: string }) => {
    // Backend expects POST /login with { email, password }
    const response = await api.post('/login', credentials);
    return response.data;
  },
  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = '/login';
  },
  getCurrentUser: async () => {
    // Backend stores user info in JWT token - decode from localStorage
    const user = localStorage.getItem('user');
    if (user) {
      return JSON.parse(user);
    }
    throw new Error('No user found');
  },
  updateProfile: async (id: string, data: any) => {
    // Use employee update endpoint
    const response = await api.put(`/employees/update/${id}`, data);
    return response.data;
  }
};
