'use client';
import { ReactNode, useState } from 'react';

interface ModalProps {
  open: boolean;
  onClose: () => void;
  title: string;
  children: ReactNode;
}

export function Modal({ open, onClose, title, children }: ModalProps) {
  if (!open) return null;
  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal slide-up" onClick={e => e.stopPropagation()}>
        <h2>{title}</h2>
        {children}
      </div>
    </div>
  );
}

interface ConfirmProps {
  open: boolean;
  onClose: () => void;
  onConfirm: () => void;
  title: string;
  message: string;
}

export function Confirm({ open, onClose, onConfirm, title, message }: ConfirmProps) {
  return (
    <Modal open={open} onClose={onClose} title={title}>
      <p style={{ color: 'var(--text-secondary)', marginBottom: 24 }}>{message}</p>
      <div className="modal-actions">
        <button className="btn btn-ghost" onClick={onClose}>Bekor qilish</button>
        <button className="btn btn-danger" onClick={() => { onConfirm(); onClose(); }}>Ha, o&apos;chirish</button>
      </div>
    </Modal>
  );
}

export function useModal() {
  const [open, setOpen] = useState(false);
  return { open, show: () => setOpen(true), hide: () => setOpen(false) };
}
