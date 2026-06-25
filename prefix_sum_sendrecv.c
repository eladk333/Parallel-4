// 322587064 - Elad Katz
#include <mpi.h>
#include <stdio.h>

int main(int argc, char** argv) {
    // Initialize the MPI environment
    MPI_Init(&argc, &argv);

    int rank, size;
    MPI_Comm_rank(MPI_COMM_WORLD, &rank);
    MPI_Comm_size(MPI_COMM_WORLD, &size);

    // Initial values per the problem definition
    int x = rank;
    int prefix = x;

    // Point-to-point communication to compute prefix sum
    if (size > 1) {
        if (rank == 0) {
            // Rank 0 just sends its prefix sum to Rank 1
            MPI_Send(&prefix, 1, MPI_INT, rank + 1, 0, MPI_COMM_WORLD);
        } else {
            int prev_prefix;
            // Receive the prefix sum from the previous rank
            MPI_Recv(&prev_prefix, 1, MPI_INT, rank - 1, 0, MPI_COMM_WORLD, MPI_STATUS_IGNORE);
            
            // Add the received sum to the current rank's x
            prefix = prev_prefix + x;
            
            // Send the updated prefix sum to the next rank, unless it's the last rank
            if (rank < size - 1) {
                MPI_Send(&prefix, 1, MPI_INT, rank + 1, 0, MPI_COMM_WORLD);
            }
        }
    }

    // Each rank prints exactly one line
    printf("rank=%d x=%d prefix=%d\n", rank, x, prefix);

    // Finalize the MPI environment
    MPI_Finalize();
    return 0;
}