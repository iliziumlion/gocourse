import api from './axios'

export const tasksAPI = {
  getAll() {
    return api.get('/tasks')
  },
  getById(id) {
    return api.get(`/tasks/${id}`)
  },
  create(task) {
    return api.post('/tasks', task)
  },
  update(id, task) {
    return api.put(`/tasks/${id}`, task)
  },
  delete(id) {
    return api.delete(`/tasks/${id}`)
  },
  getStats() {
    return api.get('/stats')
  }
}
