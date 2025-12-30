# 无需授权场景可选联系方式 - 技术设计

## 组件设计

### OptionalContactDialog.vue

**文件位置**: `frontend/src/components/export-flow/OptionalContactDialog.vue`

#### Props

```typescript
interface Props {
  visible: boolean;  // 对话框可见性
}
```

#### Emits

```typescript
interface Emits {
  (e: 'submit', data: { email?: string; additional?: string }): void;
  (e: 'skip'): void;
  (e: 'cancel'): void;
}
```

#### 内部状态

```typescript
const formData = ref({
  email: '',
  additional: '',
});
const emailError = ref('');
```

#### 邮箱校验

邮箱为可选字段，但如果填写则需要校验格式：

```typescript
const isValidEmail = (value: string) => {
  const trimmed = value.trim();
  if (!trimmed) return true; // 空值视为有效（可选）
  return emailPattern.test(trimmed);
};

const canContinue = computed(() => {
  return isValidEmail(formData.value.email);
});
```

## 状态机扩展

### 新增状态

在 `State.step` 联合类型中新增 `'optional-contact'`：

```typescript
type Step = 
  | 'idle'
  | 'confirm-signature'
  | 're-export-warning'
  | 'auth-requirement'
  | 'auth-impact-confirm'
  | 'auth-contact'
  | 'optional-contact'  // 新增
  | 'auth-gate'
  | 'auth-gate-from-picker'
  | 'auth-request'
  | 'picker'
  | 'done';
```

### 新增对话框可见性

```typescript
const optionalContactDialogVisible = ref(false);
```

### 新增处理函数

```typescript
// 用户填写联系方式后提交
const handleOptionalContactSubmit = (payload: { email?: string; additional?: string }) => {
  state.value.flowData = {
    ...(state.value.flowData ?? {}),
    contactEmail: payload.email,
    contactAdditional: payload.additional,
  };
  optionalContactDialogVisible.value = false;
  state.value.step = 'picker';
  pickerDialogVisible.value = true;
};

// 用户跳过联系方式填写
const handleOptionalContactSkip = () => {
  optionalContactDialogVisible.value = false;
  state.value.step = 'picker';
  pickerDialogVisible.value = true;
};

// 用户取消
const handleOptionalContactCancel = () => {
  optionalContactDialogVisible.value = false;
  state.value.step = 'idle';
};
```

### 流程变更

修改 `handleAuthRequirementSubmit`，使选择"无需授权"后进入可选联系方式步骤：

```typescript
const handleAuthRequirementSubmit = (payload: { requireAuthorization: boolean }) => {
  state.value.flowData = { ...(state.value.flowData ?? {}), requireAuthorization: payload.requireAuthorization };
  authRequirementDialogVisible.value = false;

  if (!payload.requireAuthorization) {
    const needSignature = state.value.flowData?.needSignature ?? true;
    if (needSignature) {
      // 变更：无需授权但需要签名 → 进入可选联系方式填写
      state.value.step = 'optional-contact';
      optionalContactDialogVisible.value = true;
    } else {
      state.value.step = 'done';
    }
    return;
  }

  // 需要授权 → 二次确认
  state.value.step = 'auth-impact-confirm';
  authImpactConfirmDialogVisible.value = true;
};
```

## UI 设计

### 对话框布局

```
┌────────────────────────────────────────┐
│ 留下联系方式（可选）                    │ 标题栏 (teal 背景)
├────────────────────────────────────────┤
│                                        │
│ 您可以选择留下联系方式，便于他人与您    │ 说明文字
│ 取得联系。此步骤为可选，可直接跳过。    │
│                                        │
│ ┌────────────────────────────────────┐ │
│ │ 邮箱                               │ │ 邮箱输入框
│ └────────────────────────────────────┘ │
│                                        │
│ ┌────────────────────────────────────┐ │
│ │ 其他联系方式（选填）                │ │ 多行文本输入
│ │                                    │ │
│ └────────────────────────────────────┘ │
│                                        │
│ 联系方式将随签名一同保存，供他人查阅    │ 提示文字
│                                        │
├────────────────────────────────────────┤
│                      [跳过] [保存并继续]│ 按钮区域
└────────────────────────────────────────┘
```

### 样式

- 标题栏背景色：`teal`（与需要授权的 `primary` 区分）
- 按钮样式：
  - 跳过：`flat`、`grey` 色
  - 保存并继续：`unelevated`、`teal` 色
