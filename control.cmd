@ECHO OFF
IF "%1%"=="" GOTO FAIL
IF "%1%"=="install" GOTO INSTALL
IF "%1%"=="uninstall" GOTO UNINSTALL
IF "%1%"=="start" GOTO START
IF "%1%"=="build" GOTO BUILD
IF "%1%"=="restart" GOTO RESTART
IF "%1%"=="stop" GOTO STOP
IF "%1%"=="status" GOTO STATUS
IF "%1%"=="clearapp" GOTO CLEARAPPS
IF "%1%"=="version" GOTO VERSION

:BUILD
::go build -o %~dp0/bin/elves-agent.exe -ldflags "-H windowsgui" %~dp0/src/agent.go
go build -o %~dp0/bin/elves-agent.exe  %~dp0/src/agent.go

%~dp0/bin/elves-agent.exe -v
GOTO EXIT

:START
net start|find /i "elves-agent"
if %errorlevel% == 0 (echo "elves-agent already started..") else (net start elves-agent)
echo "elves-agent started.."
GOTO EXIT


:INSTALL
@set username=
@set /p username=administrator username:
@set pwd=
@set /p pwd=administrator password:
%~dp0\\bin\\nssm.exe install elves-agent %~dp0\\bin\\elves-agent.exe
%~dp0\\bin\\nssm.exe set elves-agent AppParameters "-r %~dp0"
%~dp0\\bin\\nssm.exe set elves-agent Description "elves-agent"
%~dp0\\bin\\nssm.exe set elves-agent Start SERVICE_DELAYED_AUTO_START
%~dp0\\bin\\nssm.exe set elves-agent ObjectName .\%username% %pwd%
GOTO EXIT

:UNINSTALL
sc delete elves-agent
GOTO EXIT

:STARTFAIL
echo "elves-agent start fail.."
GOTO EXIT

:STOP
net start|find /i "elves-agent"
if %errorlevel% == 0 (net stop elves-agent)
echo "elves-agent stoped.."
GOTO EXIT

:STATUS
net start|find /i "elves-agent"
if %errorlevel% == 0 (echo "elves-agent started..") else (echo "elves-agent stoped..")
GOTO EXIT

:RESTART
net start|find /i "elves-agent"
if %errorlevel% == 0 (net stop elves-agent)
net start elves-agent
GOTO START

:CLEARAPPS
%~dp0/bin/elves-agent.exe -clear
GOTO EXIT

:VERSION
%~dp0/bin/elves-agent.exe -v
GOTO EXIT

:FAIL
echo "build|install|start|stop|restart|status|uninstall|cleanapp|version"

:EXIT