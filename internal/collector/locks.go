package collector

// TODO consider removing this collector
//
// func GetLocks(ctx context.Context, binPath string, repo string) {
// 	check, err := listLocks(ctx, "restic", repo)

// 	if err != nil {
// 		if errors.Is(err, ErrCheck) {
// 			slog.Error(fmt.Sprintf("Failed to get snapshot data in repo %v because: %v", repo, err))
// 		} else {
// 			metrics.CheckFailed.WithLabelValues(repo).Inc()
// 		}
// 	} else {
// 		// metrics.CheckSuccess.
// 		metrics.CheckSuccess.WithLabelValues(repo).Inc()
// 	}

// 	metrics.CheckSuggestRepairIndex.WithLabelValues(repo).Set(float64(boolToInt(check.HintRepairIndex)))
// 	metrics.CheckSuggestPrune.WithLabelValues(repo).Set(float64(boolToInt(check.HintPrune)))
// 	metrics.CheckErrorsTotal.WithLabelValues(repo).Set(float64(check.NumErrors))

// }

// func listLocks(ctx context.Context, binPath string, repo string) (*restic.CheckSummary, error) {
// 	// check and verify integrity fo the repository
// 	summary := &restic.CheckSummary{}

// 	args := []string{"-r", repo, "--no-lock", "check", "--json"}

// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		return summary, err
// 	}
// 	env := make(map[string]string)

// 	stdout := new(bytes.Buffer)
// 	stderr := new(bytes.Buffer)

// 	err = runCommand(ctx, binPath, cwd, args, env, stdout, stderr)

// 	if err == nil {
// 		if exitError, ok := err.(*exec.ExitError); ok {
// 			// cmd.Run returned with an non-zero exit code
// 			if exitError.ExitCode() == 1 {
// 				// the integrity check failed and there are probably errors
// 				return summary, ErrCheck
// 			}
// 		}
// 		return summary, err // something else happened
// 	}

// 	err = json.Unmarshal(stdout.Bytes(), summary)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return summary, nil
// }
