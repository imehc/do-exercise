import { Status } from "~/do-exercise-api";
import { Badge } from "../ui/badge";

interface Props {
    status: Status
}

export function StatusTag({ status }: Props) {
    switch (status) {
        case Status.disabled:
            return <Badge variant="outline" className="text-red-300 border-red-300">禁用</Badge>
        case Status.enabled:
            return <Badge variant="outline">启用</Badge>
        default:
            return "-"
    }
}