// 322587064 - Elad Katz
#include <mpi.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include "matrix.h"

int main(int argc, char** argv) {
    MPI_Init(&argc, &argv);

    int rank, size;
    MPI_Comm_rank(MPI_COMM_WORLD, &rank);
    MPI_Comm_size(MPI_COMM_WORLD, &size);

    if (argc != 5) {
        if (rank == 0) {
            fprintf(stderr, "Usage: mpirun -np <P> %s <N> <seedA> <seedB> <maxValue>\n", argv[0]);
        }
        MPI_Finalize();
        return 1;
    }

    // Parse command line arguments
    int N = atoi(argv[1]);
    uint64_t seedA = strtoull(argv[2], NULL, 10);
    uint64_t seedB = strtoull(argv[3], NULL, 10);
    int maxValue = atoi(argv[4]);

    // Define row distribution for current rank
    int start_row = (rank * N) / size;
    int end_row = ((rank + 1) * N) / size;
    int num_rows = end_row - start_row;
    int local_elements = num_rows * N;

    IMatrix A, B, C;
    int *sendcounts = NULL;
    int *displs = NULL;

    // Rank 0 Initialization
    if (rank == 0) {
        A = imatrix_alloc(N);
        B = imatrix_alloc(N);
        C = imatrix_alloc(N);

        imatrix_fill_random(&A, seedA, maxValue);
        imatrix_fill_random(&B, seedB, maxValue);

        // Prepare arrays for MPI_Scatterv and MPI_Gatherv
        sendcounts = (int*)malloc(size * sizeof(int));
        displs = (int*)malloc(size * sizeof(int));

        for (int i = 0; i < size; i++) {
            int r_start = (i * N) / size;
            int r_end = ((i + 1) * N) / size;
            sendcounts[i] = (r_end - r_start) * N;
            displs[i] = r_start * N;
        }
    } else {
        // Other ranks still need B allocated to receive the broadcast
        B = imatrix_alloc(N);
    }

    // 1. Send blocks of A to all ranks using MPI_Scatterv
    int *local_A = NULL;
    if (local_elements > 0) {
        local_A = (int*)malloc(local_elements * sizeof(int));
    }
    
    MPI_Scatterv(
        rank == 0 ? A.data : NULL, sendcounts, displs, MPI_INT,
        local_A, local_elements, MPI_INT,
        0, MPI_COMM_WORLD
    );

    // 2. Broadcast the entire matrix B to all ranks
    MPI_Bcast(B.data, N * N, MPI_INT, 0, MPI_COMM_WORLD);

    // 3. Compute the local block of C
    int *local_C = NULL;
    if (local_elements > 0) {
        local_C = (int*)malloc(local_elements * sizeof(int));
    }

    for (int i = 0; i < num_rows; i++) {
        for (int j = 0; j < N; j++) {
            int sum = 0;
            for (int k = 0; k < N; k++) {
                sum += local_A[i * N + k] * B.data[k * N + j];
            }
            local_C[i * N + j] = sum;
        }
    }

    // 4. Gather the computed blocks of C back to rank 0
    MPI_Gatherv(
        local_C, local_elements, MPI_INT,
        rank == 0 ? C.data : NULL, sendcounts, displs, MPI_INT,
        0, MPI_COMM_WORLD
    );

    // Output and Cleanup
    if (rank == 0) {
        long long checksum = imatrix_checksum(&C);
        printf("checksum(C)=%lld\n", checksum);

        imatrix_free(&A);
        imatrix_free(&C);
        free(sendcounts);
        free(displs);
    }
    
    // Everyone frees B and local buffers
    imatrix_free(&B);
    if (local_A) free(local_A);
    if (local_C) free(local_C);

    MPI_Finalize();
    return 0;
}