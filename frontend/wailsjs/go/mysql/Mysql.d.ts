// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {mysql} from '../models';

export function GetConsStatus(arg1:string):Promise<{[key: string]: Array<mysql.MysqlConsStatus>}>;

export function GetConspercent(arg1:string):Promise<{[key: string]: string}>;

export function GetMysqlLock(arg1:string):Promise<Array<mysql.MysqlProcessF>>;

export function GetMysqlProcesslist(arg1:string,arg2:string):Promise<Array<mysql.MysqlProcessF>>;

export function KillMysqlProcesses(arg1:string,arg2:Array<mysql.MysqlProcessF>):Promise<string>;
