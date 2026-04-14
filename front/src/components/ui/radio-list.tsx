interface RadioListOption {
    value: string | number
    label: string
}

interface RadioListProps {
    options: RadioListOption[]
    value: string | number | undefined
    onChange: (value: string | number) => void
    disabled?: boolean
}

export const RadioList = ({ options, value, onChange, disabled }: RadioListProps) => {
    return (
        <div className="space-y-2">
            {options.map((option) => {
                const isSelected = option.value === value
                return (
                    <button
                        key={option.value}
                        type="button"
                        disabled={disabled}
                        onClick={() => onChange(option.value)}
                        className={`flex w-full items-center justify-between rounded-md border px-3 py-2 text-sm hover:bg-accent disabled:cursor-not-allowed disabled:opacity-50 ${isSelected ? 'border-primary bg-primary/5 font-medium' : ''}`}
                    >
                        <span>{option.label}</span>
                        {isSelected && <span className="text-xs text-primary">&#10003;</span>}
                    </button>
                )
            })}
        </div>
    )
}
