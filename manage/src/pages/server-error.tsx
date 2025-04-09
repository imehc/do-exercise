import React from 'react'
import { Button } from '~/components/ui/button'

const ServerErrorPage: React.FC = () => {
    return (
        <div className="h-full w-full flex justify-center items-center relative overflow-hidden bg-background/80 backdrop-blur-sm">
            <div className="absolute inset-0 z-0">
                <img
                    src="/500-bg.svg"
                    alt="500 background"
                    className="w-full h-full object-cover opacity-40 animate-[float_6s_ease-in-out_infinite]"
                />
            </div>
            <div className="text-center space-y-8 relative z-10">
                <div className="relative">
                    <h1 className="text-7xl font-bold text-primary animate-[bounce_2s_ease-in-out_infinite]">
                        500
                    </h1>
                    <div className="absolute -bottom-4 left-1/2 -translate-x-1/2 w-3/4 h-1 bg-gradient-to-r from-transparent via-primary/30 to-transparent rounded-full blur-sm" />
                </div>
                <p className="text-xl text-muted-foreground">服务器出错了</p>
                <Button
                    className="px-8 py-4 text-lg relative z-20 bg-primary/90 hover:bg-primary/80 transition-colors duration-300 shadow-lg hover:shadow-xl"
                    variant="default"
                    onClick={() => window.location.reload()}
                >
                    刷新页面
                </Button>
            </div>
        </div>
    )
}

export default ServerErrorPage