package minimizer

/*
;----------------------------------------------------------------------------------------------

PRO handle_fit, ehd, event, index

qzaxis = *ehd.qzaxis
signal = *ehd.signal

IF (index EQ 99) THEN BEGIN   ; fit marked
; find fitmarks
m_all=WHERE(ehd.fitmarks,count)
If (count EQ 0) THEN m_all=[99, 99]
If (count EQ 1) THEN m_all=[m_all[0], 99]
END ELSE BEGIN  ; individual fitbutton plus dummy
m_all=[index, 99]
END
num_p=size(m_all, /N_ELEMENTS)

ehd.fitindex = m_all
Widget_Control, event.top, SET_UVALUE=ehd

parameters = adjust_params(ehd.parameters)

sub_p=DBLARR(num_p)
Xi=INTARR(num_p,num_p)
FOR i=0,num_p-1 DO BEGIN
small_Xi=INTARR(num_p)
small_Xi[i]=1
Xi[0,i]=small_Xi

IF (m_all[i] EQ 99) THEN sub_p[i]=99 ELSE sub_p[i]=parameters[m_all[i]]
END

; startwerte aus dem fit
POWELL,  sub_p,  Xi,  0.001,  Fmin,  'sim2sig_rms', /DOUBLE

FOR i=0,num_p-1 DO BEGIN
IF (m_all[i] NE 99) THEN parameters[m_all[i]]=sub_p[i]
END

ehd.parameters=adjust_params(parameters)
Widget_Control, event.top, SET_UVALUE=ehd
set_inputs, ehd
handle_simulation, ehd

RETURN
END

;----------------------------------------------------------------------------------------------
*/ //TODO understand what it does

func OlDFirParameter() {

}
