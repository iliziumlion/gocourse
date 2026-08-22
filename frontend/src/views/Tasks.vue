<script setup>
import { ref, onMounted } from 'vue'
import { tasksAPI } from '../api/tasks'

const tasks = ref([])

const loadTasks = async () => {
  try {
    const response = await tasksAPI.getAll()
    tasks.value = response.data
  } catch (error) {
    console.error('Error loading tasks:', error)
  }
}

const deleteTask = async (id) => {
  if (confirm('Вы уверены, что хотите удалить эту задачу?')) {
    try {
      await tasksAPI.delete(id)
      await loadTasks()
    } catch (error) {
      console.error('Error deleting task:', error)
    }
  }
}

const getStatusClass = (status) => {
  const classes = {
    pending: 'bg-yellow-100 text-yellow-800',
    in_progress: 'bg-blue-100 text-blue-800',
    completed: 'bg-green-100 text-green-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getStatusText = (status) => {
  const texts = {
    pending: 'В ожидании',
    in_progress: 'В работе',
    completed: 'Завершено'
  }
  return texts[status] || status
}

const getPriorityClass = (priority) => {
  const classes = {
    low: 'bg-gray-100 text-gray-800',
    medium: 'bg-yellow-100 text-yellow-800',
    high: 'bg-red-100 text-red-800'
  }
  return classes[priority] || 'bg-gray-100 text-gray-800'
}

const getPriorityText = (priority) => {
  const texts = {
    low: 'Низкий',
    medium: 'Средний',
    high: 'Высокий'
  }
  return texts[priority] || priority
}

onMounted(() => {
  loadTasks()
})
</script>

<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-3xl font-bold text-gray-800">Задачи</h2>
      <router-link
        to="/tasks/new"
        class="px-4 py-2 bg-primary text-white rounded-lg hover:bg-blue-700"
      >
        + Новая задача
      </router-link>
    </div>

    <!-- Tasks Table -->
    <div class="bg-white rounded-lg shadow overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Название</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Статус</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Приоритет</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Действия</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="task in tasks" :key="task.id">
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ task.id }}</td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-medium text-gray-900">{{ task.title }}</div>
              <div class="text-sm text-gray-500">{{ task.description }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span :class="getStatusClass(task.status)" class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                {{ getStatusText(task.status) }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span :class="getPriorityClass(task.priority)" class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                {{ getPriorityText(task.priority) }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
              <router-link :to="`/tasks/${task.id}/edit`" class="text-indigo-600 hover:text-indigo-900 mr-3">
                Редактировать
              </router-link>
              <button @click="deleteTask(task.id)" class="text-red-600 hover:text-red-900">
                Удалить
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
