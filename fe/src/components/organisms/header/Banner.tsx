import Link from "next/link";
import styles from "@/styles/Home.module.css";

export default function Banner() {
  return (
    <>
      <div className="hero">
        <div className="container">
          <div className="row align-items-center">
            <div className="col-lg-12">
              <div className="intro-wrap mt-5">
                <h1 className="mb-5 text-center poppins-bold">
                  Pesan Lapangan Futsal Impianmu Dalam Sekejap!
                </h1>
                <p className={`${styles['caption']} text-center text-white poppins-semibold`}>
                  Kumpulkan Timmu, Atur Strategi, dan Taklukkan Lapangan.
                </p>
                <div className="text-center mt-4">
                  <Link href="#field-list" className={`text-white ${styles['btn']} ${styles['btn-primary']} poppins-medium`}>
                    Booking Sekarang
                  </Link>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}