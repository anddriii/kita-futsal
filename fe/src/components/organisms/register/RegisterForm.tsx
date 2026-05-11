'use client'
import Link from "next/link";
import FormGroup from "@/components/molecules/FormGroup";
import Button from "@/components/atoms/Button";
import styles from "@/styles/Auth.module.css";
import React, {useState} from "react";
import {useRouter} from "next/navigation";
import apiConfig from "@/config/api";
import axios from "axios";
import {toast} from "react-toastify";

export default function RegisterForm() {
  const [name, setName] = useState('');
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [phoneNumber, setPhoneNumber] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [errors, setErrors] = useState<any>({});
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter()
  const handleUsernameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUsername(e.target.value);
    setFieldError('username', e.target.value);
  }

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value);
    setFieldError('password', e.target.value);
  }

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value);
    setFieldError('name', e.target.value);
  }

  const handleEmailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
    setFieldError('email', e.target.value);
  }

  const handlePhoneNumberChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPhoneNumber(e.target.value);
    setFieldError('phone_number', e.target.value);
  }

  const handleConfirmPasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setConfirmPassword(e.target.value);
    setFieldError('confirm_password', e.target.value);
  }

  const validationConditions: { [key: string]: (value: string) => boolean } = {
    username: (value: string) => value.length >= 5,
    name: (value: string) => value.length >= 3,
    email: (value: string) => value.length >= 5,
    phone_number: (value: string) => value.length >= 9,
    password: (value: string) => value.length >= 8,
    confirm_password: (value: string) => value === password
  };

  const setFieldError = (fieldName: string, fieldValue: string) => {
    if (validationConditions[fieldName](fieldValue)) {
      const newErrors = {...errors};
      delete newErrors[fieldName];
      setErrors(newErrors);
    }
  }

  const handleSubmit = async (e: React.MouseEvent<HTMLButtonElement>) => {
  e.preventDefault();
  setIsLoading(true);
  setErrors({}); // Biasakan reset error state tiap kali submit baru

  try {
    await axios.post(`${apiConfig.user.baseUrl}/api/v1/auth/register`, {
      name,
      email,
      phoneNumber,
      username,
      password,
      confirmPassword
    });

    setIsLoading(false);
    toast.success('Register berhasil');
    setTimeout(() => {
      router.push('/login');
    }, 2000);

  } catch (error: any) {
    setIsLoading(false);

    // SAFETY CHECK: Pastikan error.response beneran ada isinya
    if (error.response && error.response.data) {
      const responseData = error.response.data;
      
      // Munculin toast message dari API Golang lo
      toast.error(responseData.message || 'Terjadi kesalahan pada validasi');

      // Mapping validasi error ke form
      const newErrors: any = {};
      if (responseData.data && Array.isArray(responseData.data)) {
        responseData.data.forEach((err: any) => {
          newErrors[err.field] = err.message;
        });
        setErrors(newErrors);
      }
    } else {
      // Masuk ke sini kalau Server Down, Network Error, atau CORS Issue
      toast.error('Gagal terhubung ke server. Cek koneksi backend lo!');
      // console.error("AXIOS NETWORK ERROR:", error.message);
    }
  }
}

  return (
    <>
      <form method="post" className={`${styles['poppins-semibold']}`}>
        <div className="row">
          <div className="col-lg-6">
            <div className="form-group first">
              <FormGroup
                type="text"
                name="name"
                className={`form-control ${styles['form-input']}`}
                placeholder="Masukan Nama"
                label="Nama"
                onChange={handleNameChange}
                autoComplete={'off'}
              />
              {errors.Name ? <span className="text-xs text-danger ml-2">{errors.Name}</span> : null}
            </div>
            <div className="form-group first">
              <FormGroup
                type="text"
                name="username"
                className={`form-control ${styles['form-input']}`}
                placeholder="Masukan Username"
                label="Username"
                onChange={handleUsernameChange}
                autoComplete={'off'}
              />
              {errors.Username ? <span className="text-xs text-danger ml-2">{errors.Username}</span> : null}
            </div>
            <div className="form-group first">
              <FormGroup
                type="email"
                name="email"
                className={`form-control ${styles['form-input']}`}
                placeholder="Masukan Email"
                label="Email"
                onChange={handleEmailChange}
                autoComplete={'off'}
              />
              {errors.Email ? <span className="text-xs text-danger ml-2">{errors.Email}</span> : null}
            </div>
          </div>
          <div className="col-lg-6">
            <div className="form-group last mb-3">
              <FormGroup
                type="text"
                name="phone_number"
                className={`form-control ${styles['form-input']}`}
                placeholder="Masukan No Hp"
                label="Nomor Hp."
                onChange={handlePhoneNumberChange}
                autoComplete={'off'}
              />
              {errors.PhoneNumber ? <span className="text-xs text-danger ml-2">{errors.PhoneNumber}</span> : null}
            </div>
            <div className="form-group last mb-3">
              <FormGroup
                type="password"
                name="password"
                className={`form-control ${styles['form-input']}`}
                placeholder="Masukan Password"
                label="Password"
                onChange={handlePasswordChange}
              />
              {errors.Password ? <span className="text-xs text-danger ml-2">{errors.Password}</span> : null}
            </div>
            <div className="form-group last mb-3">
              <FormGroup
                type="password"
                name="confirm_password"
                className={`form-control ${styles['form-input']}`}
                placeholder="Masukan Konfirm Password"
                label="Konfrimasi Password"
                onChange={handleConfirmPasswordChange}
              />
              {errors.ConfirmPassword ?
                <span className="text-xs text-danger ml-2">{errors.ConfirmPassword}</span> : null}
            </div>
          </div>
        </div>
        <Button
          disabled={isLoading}
          type="button"
          className={`btn btn-block ${styles['btn-register']}`}
          onClick={handleSubmit}
        >
          {isLoading ? 'Loading...' : 'Register'}
        </Button>
        <div className="d-flex mb-5 align-items-center mt-1">
          <span className="ml-auto">
            <Link
              href="/login"
              className={`${styles['forgot-pass']} ${styles['poppins-semibold']}`}
              style={{textDecoration: 'none'}}
            >
              Sudah punya akun?
              <strong style={{textDecoration: 'underline', marginLeft: '5px'}}>
                Login disini
              </strong>
            </Link>
          </span>
        </div>
      </form>
    </>
  )
}