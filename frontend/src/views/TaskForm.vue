<template>
  <div>
    <h2 class="text-3xl font-bold text-gray-800 mb-6">
      {{ isEdit ? 'Редактировать задачу' : 'Новая задача' }}
    </h2>

    <div class="bg-white rounded-lg shadow p-6 max-w-2xl">
      <form @submit.prevent="saveTask">
        <div class="mb-4">
          <label class="block text-gray-700 text-sm font-bold mb-2">Название</label>
          <input
            v-model="form.title"
            type="text"
            required
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            placeholder="Введите название задачи"
          >
        </div>

        <div class="mb-4">
          <label class="block text-gray-700 text-sm font-bold mb-2">Описание</label>
          <textarea
            v-model="form.description"
            rows="4"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            placeholder="Введите описание задачи"
          ></textarea>
        </div>

        <div class="mb-4">
          <label class="block text-gray-700 text-sm font-bold mb-2">Статус</label>
          <select
            v-model="form.status"
            class="shadow border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          >
            <option value="pending">В ожидании</option>
            <option value="in_progress">В работе</option>
            <option value="completed">Завершено</option>
          </select>
        </div>

        <div class="mb-6">
          <label class="block text-gray-700 text-sm font-bold mb-2">Приоритет</label>
          <select
            v-model="form.priority"
            class="shadow border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          >
            <option value="low">Низкий</option>
            <option value="medium">Средний</option>
            <option value="high">Высокий</option>
          </select>
        </div>

        <div class="flex items-center justify-between">
          <button
            type="submit"
            class="bg-primary hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
          >
            {{ isEdit ? 'Обновить' : 'Создать' }}
          </button>
          <router-link
            to="/tasks"
            class="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800"
          >
            Отмена
          </router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { tasksAPI } from '../api/tasks'

const route = useRoute()
const router = useRouter()

const form = ref({
  title: '',
  description: '',
  status: 'pending',
  priority: 'medium'
})

const isEdit = computed(() => !!route.params.id)

const loadTask = async () => {
  if (isEdit.value) {
    try {
      const response = await tasksAPI.getById(route.params.id)
      form.value = response.data
    } catch (error) {
      console.error('Error loading task:', error)
    }
  }
}

const saveTask = async () => {
  try {
    if (isEdit.value) {
      await tasksAPI.update(route.params.id, form.value)
    } else {
      await tasksAPI.create(form.value)
    }
    router.push('/tasks')
  } catch (error) {
    console.error('Error saving task:', error)
  }
}

onMounted(() => {
  loadTask()
})
</script>
