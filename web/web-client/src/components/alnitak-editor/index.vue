<template>
  <div class="alnitak-editor">
    <WangToolbar class="toolbar" :editor="editorRef" :defaultConfig="toolbarConfig" :mode="mode" />
    <WangEditor class="editor" v-model="valueHtml" :defaultConfig="editorConfig" :mode="mode" @onCreated="handleCreated"
      @onChange="handleChange" @onDestroyed="handleDestroyed" @onFocus="handleFocus" @onBlur="handleBlur"
      @customAlert="customAlert" @customPaste="customPaste" />
  </div>
</template>

<script setup lang="ts">
import { uploadFileAPI } from '@/api/upload';
import '@wangeditor/editor/dist/css/style.css';
import { onBeforeUnmount, ref, shallowRef, onMounted } from 'vue';
import { ElMessage } from 'element-plus';

const emit = defineEmits(['update:content']);
const props = withDefaults(defineProps<{
  content: string;
}>(), {
  content: "",
});

const mode = 'default';
const editorRef = shallowRef();
const valueHtml = ref(props.content);

const toolbarConfig = {
  excludeKeys: [
    "blockquote", "bgColor", "fontSize", "fontFamily", "lineHeight", "todo",
    "group-indent", "emotion", "insertTable", "group-image", "group-video",
    "codeBlock", "fullScreen"
  ],
  insertKeys: {
    index: 20,
    keys: "uploadImage"
  }
};

const editorConfig = {
  height: 200,
  placeholder: '请输入文章内容...',
  MENU_CONF: {
    uploadImage: {
      async customUpload(file: File, insertFn: any) {
        try {
          const result = await uploadFileAPI({
            name: "image",
            action: "v1/upload/image",
            file,
            onProgress: () => {},
            onError: () => ElMessage.error("上传失败"),
            onFinish: (data: any) => insertFn(getResourceUrl(data.data.url), "图片", "")
          });
        } catch (error) {
          ElMessage.error("上传失败");
        }
      }
    }
  }
};

const handleCreated = (editor: any) => {
  editorRef.value = editor;
};

const handleChange = (editor: any) => {
  emit("update:content", editor.getHtml());
};

const handleDestroyed = (editor: any) => {
  console.log('destroyed', editor);
};

const handleFocus = (editor: any) => {
  console.log('focus', editor);
};

const handleBlur = (editor: any) => {
  console.log('blur', editor);
};

const customAlert = (info: any, type: any) => {
  alert(`【自定义提示】${type} - ${info}`);
};

const customPaste = (editor: any, event: any, callback: any) => {
  console.log('ClipboardEvent 粘贴事件对象', event);
  const text = event.clipboardData.getData('text/plain');
  editor.insertText(text);
  event.preventDefault();
  callback(false);
};

onBeforeUnmount(() => {
  const editor = editorRef.value;
  if (editor) editor.destroy();
});
</script>

<style lang="scss" scoped>
.alnitak-editor {
  border: 1px solid #ccc;

  .toolbar {
    border-bottom: 1px solid #ccc;
  }
}

:deep(.w-e-text-container) {
  min-height: 300px !important;
}
</style>
