/**
 * Typed Redux Hooks
 * Pre-typed versions of useDispatch and useSelector
 * Use these throughout your app instead of plain `useDispatch` and `useSelector`
 */

import { TypedUseSelectorHook, useDispatch, useSelector } from 'react-redux';
import type { AppDispatch, RootState } from './store';

// Use throughout your app instead of plain `useDispatch` and `useSelector`
export const useAppDispatch = () => useDispatch<AppDispatch>();
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;

