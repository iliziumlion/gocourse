import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import Tasks from '../views/Tasks.vue'
import TaskForm from '../views/TaskForm.vue'

const routes = [
  { path: '/', name: 'Dashboard', component: Dashboard },
  { path: '/tasks', name: 'Tasks', component: Tasks },
  { path: '/tasks/new', name: 'NewTask', component: TaskForm },
  { path: '/tasks/:id/edit', name: 'EditTask', component: TaskForm }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
