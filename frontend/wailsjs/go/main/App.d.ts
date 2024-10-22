// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {configuration} from '../models';
import {database} from '../models';

export function AddConfiguration(arg1:configuration.DatabaseConnection):Promise<void>;

export function GetAvailableTables():Promise<any>;

export function GetConfiguration(arg1:string):Promise<configuration.DatabaseConnection>;

export function GetCurrentConfigurations():Promise<Array<string>>;

export function GetTableData(arg1:string,arg2:string):Promise<database.TableData>;

export function SetActiveConnection(arg1:number):Promise<void>;

export function TestConnection(arg1:configuration.DatabaseConnection):Promise<void>;

export function UpdateValue(arg1:database.UpdateValueRequest):Promise<void>;
