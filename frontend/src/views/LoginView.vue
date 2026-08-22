<template>
  <div class="auth-card">
    <h1>go-waste-routes</h1>
    <p>本地管理员：admin / Admin123!</p>
    <el-form @submit.prevent="submit">
      <el-form-item><el-input v-model="form.username" placeholder="用户名" /></el-form-item>
      <el-form-item><el-input v-model="form.password" type="password" placeholder="密码" show-password /></el-form-item>
      <el-button type="primary" native-type="submit" :loading="loading">登录</el-button>
      <el-button text @click="$router.push('/register')">注册</el-button>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const loading = ref(false)
const form = reactive({ username: 'admin', password: 'Admin123!' })

async function submit() {
  loading.value = true
  try {
    await auth.login(form)
    router.push('/')
  } finally {
    loading.value = false
  }
}
</script>
