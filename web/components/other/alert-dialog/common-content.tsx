import {
  AlertDialogHeader,
  AlertDialogFooter,
  AlertDialogContent,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogCancel,
  AlertDialogAction,
} from '~/components/ui/alert-dialog';

interface Props {
  onCancel(): void;
  onContinue(): void;
  title: React.ReactNode;
  subTitle?: React.ReactNode;
}

// TODO: 多语言
export function CommonAlertDialogContent({ onCancel, onContinue, title, subTitle }: Props) {
  return (
    <AlertDialogContent>
      <AlertDialogHeader>
        <AlertDialogTitle>{title}</AlertDialogTitle>
        <AlertDialogDescription>{subTitle}</AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel onClick={onCancel}>取消</AlertDialogCancel>
        <AlertDialogAction onClick={onContinue}>确定</AlertDialogAction>
      </AlertDialogFooter>
    </AlertDialogContent>
  );
}
