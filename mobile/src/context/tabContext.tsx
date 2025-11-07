import React, {createContext} from 'react';

type TabContextType = {
  year: number;
  month: number;
  setYear: (year: number) => void;
  setMonth: (month: number) => void;
};

export const TabContext = createContext<TabContextType | undefined>(undefined);

export const TabProvider = TabContext.Provider;

export const useTab = () => {
  const context = React.useContext(TabContext);
  if (context === undefined) {
    throw new Error('useTab must be used within a TabProvider');
  }
  return context;
};
