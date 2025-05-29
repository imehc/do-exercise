import { describe, it, expect } from 'vitest'
import { schema } from './action-schema'

describe('User Action Schema Validation', () => {
  // Helper to simplify validation
  const validate = (input: unknown) => schema.safeParse(input)

  it('should pass when isEdit is true and password fields are empty', () => {
    const result = validate({
      isEdit: true,
      password: '',
      confirmPassword: '',
      username: 'testuser',
      roleIds: [1],
    })

    expect(result.success).toBe(true)
  })

  it('should fail when isEdit is false and password is empty', () => {
    const result = validate({
      isEdit: false,
      password: '',
      confirmPassword: '',
      username: 'testuser',
      roleIds: [1],
    })

    expect(result.success).toBe(false)
    if (!result.success) {
      expect(result.error.issues).toContainEqual(
        expect.objectContaining({
          path: ['password'],
          message: '请输入密码',
        })
      )
      expect(result.error.issues).toContainEqual(
        expect.objectContaining({
          path: ['confirmPassword'],
          message: '请输入确认密码',
        })
      )
    }
  })

  it('should fail when password is too short', () => {
    const result = validate({
      isEdit: false,
      password: 'Abc1!',
      confirmPassword: 'Abc1!',
      username: 'testuser',
      roleIds: [1],
    })

    expect(result.success).toBe(false)
    if (!result.success) {
      expect(result.error.issues).toContainEqual(
        expect.objectContaining({
          message: '密码至少为6个字符',
        })
      )
    }
  })

  it('should fail when password lacks a letter', () => {
    const result = validate({
      isEdit: false,
      password: '123456!',
      confirmPassword: '123456!',
      username: 'testuser',
      roleIds: [1],
    })

    expect(result.success).toBe(false)
    if (!result.success) {
      expect(result.error.issues).toContainEqual(
        expect.objectContaining({
          message: '密码必须包含至少一个字母',
        })
      )
    }
  })

  it('should fail when password lacks a number', () => {
    const result = validate({
      isEdit: false,
      password: 'Abcdefg!',
      confirmPassword: 'Abcdefg!',
      username: 'testuser',
      roleIds: [1],
    })

    expect(result.success).toBe(false)
    if (!result.success) {
      expect(result.error.issues).toContainEqual(
        expect.objectContaining({
          message: '密码必须包含至少一个数字',
        })
      )
    }
  })

  it('should fail when password lacks a special character', () => {
    const result = validate({
      isEdit: false,
      password: 'Abcdefg1',
      confirmPassword: 'Abcdefg1',
      username: 'testuser',
      roleIds: [1],
    })

    expect(result.success).toBe(false)
    if (!result.success) {
      expect(result.error.issues).toContainEqual(
        expect.objectContaining({
          message: '密码必须包含至少一个特殊字符',
        })
      )
    }
  })

  it('should succeed with valid password and confirm password', () => {
    const result = validate({
      isEdit: false,
      password: 'Abcdefg1!',
      confirmPassword: 'Abcdefg1!',
      username: 'testuser',
      roleIds: [1],
    })

    expect(result.success).toBe(true)
  })

  it('should fail when passwords do not match', () => {
    const result = validate({
      isEdit: false,
      password: 'Abcdefg1!',
      confirmPassword: 'Xyz1234!',
      username: 'testuser',
      roleIds: [1],
    })

    expect(result.success).toBe(false)
    if (!result.success) {
      expect(result.error.issues).toContainEqual(
        expect.objectContaining({
          path: ['confirmPassword'],
          message: '两次输入的密码不一致',
        })
      )
    }
  })

  it('should pass when all required fields are provided', () => {
    const result = validate({
      isEdit: false,
      password: 'Abcdefg1!',
      confirmPassword: 'Abcdefg1!',
      username: 'testuser',
      roleIds: [1],
      email: 'test@example.com',
    })

    expect(result.success).toBe(true)
  })
})
