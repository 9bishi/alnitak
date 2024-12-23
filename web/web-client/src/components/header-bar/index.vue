<template>
  <div class="header-bar">
    <nuxt-link class="header-left" to="/">
      <h1 class="title">{{ globalConfig.title }}</h1>
    </nuxt-link>
    
    <!-- 搜索框 -->
    <div v-if="isSearchPage" class="header-center">
      <div class="search-form">
        <input 
          class="input" 
          v-model="keywords" 
          @keydown.enter="debouncedSearch"
          placeholder="请输入搜索内容" 
        />
        <button class="btn" @click="debouncedSearch">
          <search-icon class="icon" size="16" />
        </button>
      </div>
    </div>

    <div v-if="!loading" class="header-right">
      <!-- 用户头像和菜单 -->
      <div v-if="isLoggedIn" class="avatar-box">
        <nuxt-link to="/space">
          <common-avatar :url="userInfo?.avatar" :size="40" :iconSize="22"></common-avatar>
        </nuxt-link>
        <div class="menu-container">
          <div class="transition"></div>
          <div class="header-menu">
            <div class="menu-info">
              <div class="name-box">
                <span class="name">{{ userInfo?.name }}</span>
                <span class="sign">{{ userInfo?.sign }}</span>
              </div>
            </div>
            <div class="divider disabled-select"></div>
            <nuxt-link class="menu-item disabled-select" to="/space">
              <div class="link-title">
                <user-icon class="icon" />
                <span>个人中心</span>
              </div>
              <right-icon class="right-icon" />
            </nuxt-link>
            <div class="menu-item disabled-select" @click="logout">
              <div class="link-title">
                <logout-icon class="icon" />
                <span>退出登录</span>
              </div>
              <right-icon class="right-icon" />
            </div>
          </div>
        </div>
      </div>
      
      <div v-else class="avatar-box">
        <nuxt-link class="login-btn" to="/login">登录</nuxt-link>
      </div>

      <!-- 图标按钮 -->
      <nuxt-link class="icon-btn" to="/message/announce">
        <message-icon class="icon" />
        <div class="icon-text">消息</div>
      </nuxt-link>
      <nuxt-link class="icon-btn" to="/space/history">
        <history-icon class="icon" />
        <div class="icon-text">历史</div>
      </nuxt-link>
      <nuxt-link class="icon-btn" to="/space/collection">
        <collect-icon class="icon" />
        <div class="icon-text">收藏</div>
      </nuxt-link>
      
      <!-- 投稿按钮 -->
      <nuxt-link class="upload-btn disabled-select" to="/upload/video">
        <upload-icon class="upload-icon" />
        <span>投稿</span>
      </nuxt-link>
    </div>
    
    <div v-else class="header-right"></div>
  </div>
</template>
<script setup lang="ts">
import { ref, onBeforeMount } from 'vue';
import { debounce } from 'lodash';
import { useRoute, navigateTo } from 'vue-router';
import { ElMessage } from 'element-plus';
import { logoutAPI } from '@/api/auth';
import { getUserInfoAPI } from '@/api/user';
import CommonAvatar from '@/components/common-avatar/index.vue';
import {
  Search as SearchIcon, Message as MessageIcon,
  Upload as UploadIcon, History as HistoryIcon,
  User as UserIcon, Logout as LogoutIcon,
  Right as RightIcon, FolderFocusOne as CollectIcon
} from '@icon-park/vue-next';

// 全局配置与路由
const route = useRoute();
const isSearchPage = ref(route.name !== 'search-keywords');
const globalConfig = { title: '您的网站' };  // 假设的全局配置

// 搜索功能
const keywords = ref('');
const debouncedSearch = debounce(() => {
  if (!keywords.value) {
    ElMessage.warning("请先输入搜索内容");
    return;
  }
  navigateTo(`/search/${keywords.value}`, { open: { target: '_blank' } });
}, 500);

// 用户信息与登录状态
const loading = ref(true);
const isLoggedIn = ref(false);
const userInfo = ref<any>(null);  // 更灵活的用户信息类型

const getUserInfo = async () => {
  try {
    const res = await getUserInfoAPI();
    if (res.data.code === 200) {  // 假设 `statusCode.OK` 为 200
      userInfo.value = res.data.data.userInfo;
      isLoggedIn.value = true;
    }
  } catch (error) {
    console.error("获取用户信息失败", error);
    ElMessage.error("加载用户信息失败");
  } finally {
    loading.value = false;
  }
};

// 退出登录
const logout = async () => {
  try {
    await logoutAPI(localStorage.getItem('refreshToken'));
    localStorage.removeItem("token");
    localStorage.removeItem('refreshToken');
    isLoggedIn.value = false;
  } catch (error) {
    console.error("退出登录失败", error);
    ElMessage.error("退出登录失败");
  }
};

// 生命周期钩子
onBeforeMount(() => {
  getUserInfo();
});
</script>
<style lang="scss" scoped>
.header-bar {
  width: 100%;
  height: 60px;
  padding: 0 36px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: #fff;
  box-shadow: 0px 0px 3px #c8c8c8;
  z-index: 999;

  .header-left {
    display: flex;
    align-items: center;
    width: 260px;
    .title {
      font-size: 16px;
      font-weight: 500;
      cursor: pointer;
      overflow: hidden;
      white-space: nowrap;
    }
  }

  .header-center {
    width: 300px;
    .search-form {
      position: relative;
      .input {
        padding: 8px 30px 8px 11px;
        height: 36px;
        font-size: 12px;
        border-radius: 18px;
        width: 300px;
      }
      .btn {
        position: absolute;
        top: 0;
        right: 10px;
        width: 20px;
        height: 36px;
        background: transparent;
        cursor: pointer;
      }
    }
  }

  .header-right {
    width: 320px;
    display: flex;
    align-items: center;
    .avatar-box {
      margin-right: 10px;
    }
    .login-btn {
      display: block;
      width: 40px;
      height: 40px;
      border-radius: 50%;
      background-color: var(--primary-hover-color);
      text-align: center;
      line-height: 40px;
      font-size: 14px;
    }
    .upload-btn {
      display: flex;
      align-items: center;
      background-color: var(--primary-color);
      width: 100px;
      height: 36px;
      border-radius: 5px;
      text-align: center;
      margin-left: 10px;
      cursor: pointer;
    }
  }
}

.icon-btn {
  color: var(--font-primary-3);
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  .icon {
    width: 18px;
    height: 18px;
    margin-bottom: 2px;
  }
  .icon-text {
    font-size: 12px;
  }
}

.avatar-box {
  position: relative;
  cursor: pointer;
  &:hover .menu-container {
    display: block;
  }
  .menu-container {
    display: none;
    position: absolute;
    top: 40px;
    left: -110px;
    width: 230px;
    box-shadow: 0 0 30px rgba(0, 0, 0, .1);
    background-color: #fff;
    padding: 12px;
    border-radius: 10px;
    .menu-item {
      display: flex;
      align-items: center;
      height: 40px;
      padding-left: 10px;
      .link-title {
        display: flex;
        align-items: center;
        .icon {
          width: 18px;
          height: 18px;
          margin-right: 10px;
        }
        span {
          font-size: 14px;
        }
      }
    }
    .divider {
      height: 1px;
      background-color: #e0e0e0;
      margin: 8px 0;
    }
  }
}
</style>
