import React, {Dispatch, SetStateAction, useState} from 'react';
import {StyleSheet, View} from 'react-native';
import DropDownPicker from 'react-native-dropdown-picker';

import colors from '../../../colors';

export type DropDownYearProps = {
  year: number;
  years: number[];
  setYear: Dispatch<SetStateAction<number>>;
  onSelected?: (year: number) => void;
};

const styles = StyleSheet.create({
  picker: {
    borderWidth: 0,
    backgroundColor: colors.main_bg,
    zIndex: 1000,
    marginBottom: 0,
    width: 60,
    padding: 0,
    paddingVertical: 0,
    paddingHorizontal: 0,
    margin: 0,

    minHeight: 0,
    justifyContent: 'center',
    gap: 0,
    // position: 'absolute',
  },
  container: {
    borderWidth: 0,
    padding: 0,
    margin: 0,
  },
  dropdownContainer: {
    borderWidth: 0.5,
    borderRadius: 5,
    width: 100,
    borderColor: colors.light_gray,
    shadowOffset: {width: 0, height: 2},
  },
  label: {
    fontWeight: 'bold',
  },
  arrowIcon: {
    width: 7,
  },
  arrowIconContainer: {
    width: 'auto',
    padding: 0,
    margin: 0,
  },
  text: {
    margin: 0,
    padding: 0,
    fontSize: 14,
    color: colors.dark_gray,
  },
});

export const DropdownYear = ({
  year,
  years,
  setYear,
  onSelected,
}: DropDownYearProps) => {
  const [open, setOpen] = useState(false);

  return (
    <View className="my-1 w-1/3 flex items-center flex-row justify-center">
      <View>
        <DropDownPicker
          style={styles.picker}
          open={open}
          value={year}
          items={
            years.length !== 0
              ? years.map(y => ({label: y + '', value: y}))
              : [
                  {
                    label: new Date().getFullYear().toString(),
                    value: new Date().getFullYear(),
                  },
                ]
          }
          setOpen={setOpen}
          setValue={setYear}
          onSelectItem={
            ((item: {label: string; value: number}) => {
              onSelected && onSelected(item.value);
              return;
            }) as any
          }
          containerStyle={styles.container}
          labelStyle={styles.label}
          dropDownContainerStyle={styles.dropdownContainer}
          arrowIconStyle={styles.arrowIcon}
          arrowIconContainerStyle={styles.arrowIconContainer}
          textStyle={styles.text}
        />
      </View>
    </View>
  );
};
